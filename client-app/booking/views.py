from django.http import JsonResponse
from django.views import View
from django.urls import reverse
from django.contrib.auth.decorators import login_required
from django.shortcuts import render
import json
from .models import BookingTransaction

@login_required
def list_of_bookings(request):
    bookings = json.loads(list_bookings(request).content)['user_bookings']
    return render(request, 'booking/bookings.html', {'bookings': bookings})

# Create your views here.
@login_required
def list_bookings(request):
    user_bookings = BookingTransaction.objects.filter(user=request.user)
    status_filter = request.GET.get('status')
    if status_filter:
        user_bookings = user_bookings.filter(status=status_filter)
    
    bookings_list = list(user_bookings.values())
    return JsonResponse({'user_bookings': bookings_list})

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
        return JsonResponse({"message": "Booking created successfully", "booking_id": booking.uuid})

@login_required
def refresh_booking_status(request, booking_id):
    try:
        booking = BookingTransaction.objects.get(uuid=booking_id)
        
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