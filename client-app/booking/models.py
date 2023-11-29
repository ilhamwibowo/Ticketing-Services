import uuid
from django.contrib.auth.models import User
from django.db import models
import requests
import json
class TicketingExternalAPI():
    url = "http://app:3000"

    def get_events():
        endpoint = "events"
        response = requests.get(f"{TicketingExternalAPI.url}/{endpoint}")
        events = json.loads(response.text)
        print(events)
        return events

    def get_seats(event_id):
        endpoint = f"events/{event_id}/empty-seats"
        response = requests.get(f"{TicketingExternalAPI.url}/{endpoint}")
        events = json.loads(response.text)
        print(events)
        return events
    
    def get_seat_status(seat_id):
        pass

    def hold_seat(event_id, seat_number):
        endpoint = f"book/{event_id}/{seat_number}"
        response = requests.post(f"{TicketingExternalAPI.url}/{endpoint}")
        events = json.loads(response.text)
        print(events)
        return events

class Invoice(models.Model):
    id = models.CharField(primary_key=True, max_length=100, default=None)  # ID as a string
    invoice = models.FileField()  # Store the PDF file
    status = models.CharField(max_length=100, default=None, null=True)
    payment_url = models.URLField(default=None, null=True)

    def save(self, *args, **kwargs):
        # Generate UUID if not provided
        if not self.id:
            self.id = uuid.uuid4()
        super().save(*args, **kwargs)

class BookingTransaction(models.Model):
    id = models.UUIDField(primary_key=True, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    event_id = models.CharField(max_length=100, default=None)
    seats = models.JSONField(default=list)  # List of booked seats
    invoice = models.OneToOneField(Invoice, on_delete=models.CASCADE, null=True)

    def save(self, *args, **kwargs):
        # Generate UUID if not provided
        if not self.id:
            self.id = uuid.uuid4()
        super().save(*args, **kwargs)
    
    def get_all_bookings_for_user(user: User):
        return BookingTransaction.objects.filter(user=user)