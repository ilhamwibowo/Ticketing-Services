from django.conf import settings
from django.http import JsonResponse, FileResponse
from django.views import View
from django.contrib.auth.decorators import login_required
from django.shortcuts import render, get_object_or_404
import json
from .models import BookingTransaction, Invoice, TicketingExternalAPI
from django.db.models import F
from django.views.decorators.csrf import csrf_exempt


@login_required
def list_of_bookings(request):
    bookings = BookingTransaction.get_all_bookings_for_user(request.user)
    print(bookings)
    return render(request, 'booking/bookings.html', {'bookings': bookings})

@login_required
def list_bookings(request):
    user_bookings = BookingTransaction.get_all_bookings_for_user(request.user)
    status_filter = request.GET.get('status')
    if status_filter:
        user_bookings = user_bookings.filter(status=status_filter)
    return JsonResponse({'user_bookings': list(user_bookings)})

class BookView(View):
    template_name = 'booking/book.html'  # HTML template file
    
    def get(self, request):
        events = TicketingExternalAPI.get_events()
        return render(request, self.template_name, {'events': events})

    def post(self, request):
        # Extract data from the request
        request_body = json.loads(request.body)
        event_id = request_body.get('event_id')
        seats = request_body.get('seats')
        status = 'PENDING'
        user = request.user  # Current logged-in user

        for seat in seats:
            # Create a new invoice
            hold_seat_req = TicketingExternalAPI.hold_seat(seat)
            invoice = Invoice(
                id=hold_seat_req['invoice_id']
            )
            invoice.save()
        
            # Create and save the BookingTransaction entry
            booking = BookingTransaction.objects.create(
                event_id=event_id,
                seats=[seat],
                status=status,
                user=user,
                invoice=invoice
            )
        
        # Return a success message or the created entry data
        return JsonResponse({"message": "Booking created successfully", "booking_id": booking.id})

@login_required
def refresh_booking_status(request, booking_id):
    try:
        booking = BookingTransaction.objects.get(id=booking_id)
        
        # Simulating external API call and data retrieval
        # Replace this logic with your actual external API call to update the status
        # For example, assuming 'new_status' is received from the external API
        new_status = 'SUCCESS'  # Replace this with the actual status received from the API

        # Update the booking status and save it
        booking.status = new_status
        booking.save()

        return JsonResponse({'message': f'Booking {booking_id} status updated to {new_status}'})
    
    except BookingTransaction.DoesNotExist:
        return JsonResponse({'error': 'Booking not found'}, status=404)
    
def get_available_events(_):
    available_events = TicketingExternalAPI.get_events()
    return JsonResponse({'available_events': available_events})

def get_chairs_status(_, event_id):
    chairs_status = TicketingExternalAPI.get_seats(event_id)
    return JsonResponse({'chairs_status': chairs_status})

@csrf_exempt
def create_invoice(request):
    if request.method == 'POST':
        invoice_id = request.POST.get('invoice_id')
        pdf_invoice = request.FILES.get('invoice_pdf')

        if not pdf_invoice:
            return JsonResponse({'message': 'No Files'})

        try:
            invoice = Invoice.objects.get(id=invoice_id)
            invoice.invoice = pdf_invoice
            invoice.save()
            return JsonResponse({
                'message': 'Invoice updated successfully',
                'invoice_id': invoice_id
            })
        except Invoice.DoesNotExist:
            # If the invoice doesn't exist, create a new one
            invoice = Invoice(id=invoice_id, invoice=pdf_invoice)
            invoice.save()
            return JsonResponse({
                'message': 'Invoice created successfully',
                'invoice_id': invoice_id
            })
    
    return JsonResponse({'error': 'Invalid request method'})


@csrf_exempt
def get_invoice(_, invoice_id):
    invoice = Invoice.objects.get(id=invoice_id)
    if invoice.invoice:  
        return FileResponse(invoice.invoice)
    else:
        return JsonResponse({ 'message': 'No File' })