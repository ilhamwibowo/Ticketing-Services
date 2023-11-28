from django.conf import settings
from django.http import JsonResponse, FileResponse
from django.views import View
from django.contrib.auth.decorators import login_required
from django.shortcuts import render, get_object_or_404
import json
from .models import BookingTransaction, Invoice
from django.db.models import F
from django.views.decorators.csrf import csrf_exempt


@login_required
def list_of_bookings(request):
    bookings = json.loads(list_bookings(request).content)['user_bookings']
    print(bookings)
    return render(request, 'booking/bookings.html', {'bookings': bookings})

@login_required
def list_bookings(request):
    user_bookings = BookingTransaction.objects.filter(user=request.user)
    status_filter = request.GET.get('status')
    if status_filter:
        user_bookings = user_bookings.filter(status=status_filter)
    
    # Perform left join with Invoice model
    bookings_with_invoices = user_bookings.annotate(
        invoice_id=F('invoice__id'),
        invoice_transaction_id=F('invoice__transaction_id'),
        invoice_file=F('invoice__invoice')
    ).values(
        'id',  # Select fields from BookingTransaction
        'event_id',
        'seats',
        'status',
        'invoice_id',  # Fields from Invoice model
        'invoice_transaction_id',
        'invoice_file'  # File field
    )

    return JsonResponse({'user_bookings': list(bookings_with_invoices)})

class BookView(View):
    template_name = 'booking/book.html'  # HTML template file
    
    def get(self, request):
        # Retrieve available events using Django URL name
        events_response = get_available_events(request).content
        events = json.loads(events_response)['available_events']

        return render(request, self.template_name, {'events': events})

    def post(self, request):
        # Extract data from the request
        request_body = json.loads(request.body)
        event_id = request_body.get('event_id')
        seats_booked = request_body.get('seats')
        status = 'PENDING'
        user = request.user  # Current logged-in user
        
        # Create and save the BookingTransaction entry
        booking = BookingTransaction.objects.create(
            event_id=event_id,
            seats=seats_booked,
            status=status,
            user=user
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
    
def get_available_events(request):
    available_events = [{'id': '1', 'name': 'Event 1'}, {'id': '2', 'name': 'Event 2'}]
    return JsonResponse({'available_events': available_events})

def get_chairs_status(request, event_id):
    chairs_status = {'event_id': event_id, 'chairs': ['A1', 'A2', 'B1', 'B2']}
    return JsonResponse({'chairs_status': chairs_status})

@csrf_exempt
def create_invoice(request):
    if request.method == 'POST':
        booking_id = request.POST.get('booking_id')
        pdf_invoice = request.FILES.get('pdf_invoice')

        if not pdf_invoice:
            return JsonResponse({'message': 'No Files'})

        # Retrieve the booking transaction object
        booking_transaction = get_object_or_404(BookingTransaction, id=booking_id)

        # Create an Invoice object
        new_invoice = Invoice(transaction=booking_transaction, invoice=pdf_invoice)
        new_invoice.save()

        return JsonResponse({'message': 'Invoice created successfully'})
    
    return JsonResponse({'error': 'Invalid request method'})


@csrf_exempt
def get_invoice(_, invoice_id):
    invoice = Invoice.objects.get(id=invoice_id)  
    return FileResponse(invoice.invoice)