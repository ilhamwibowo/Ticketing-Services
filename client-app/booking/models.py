import uuid
from django.contrib.auth.models import User
from django.db import models

class TicketingExternalAPI():
    def get_events():
        return [{'id': '1', 'name': 'Event 1'}, {'id': '2', 'name': 'Event 2'}]

    def get_seats(event_id):
        return {'event_id': event_id, 'chairs': ['A1', 'A2', 'B1', 'B2']}
    
    def get_seat_status(seat_id):
        pass

    def hold_seat(seat_id):
        return {
            'invoice_id': uuid.uuid4()
        }

class Invoice(models.Model):
    id = models.CharField(primary_key=True, max_length=100, default=None)  # ID as a string
    invoice = models.FileField()  # Store the PDF file

class BookingTransaction(models.Model):
    STATUS_CHOICES = (
        ('SUCCESS', 'Success'),
        ('FAILED', 'Failed'),
        ('PENDING', 'Pending'),
    )
    id = models.UUIDField(primary_key=True, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    event_id = models.CharField(max_length=100, default=None)
    seats = models.JSONField(default=list)  # List of booked seats
    status = models.CharField(max_length=10, choices=STATUS_CHOICES, default=None)
    invoice = models.OneToOneField(Invoice, on_delete=models.CASCADE, null=True)

    def save(self, *args, **kwargs):
        # Generate UUID if not provided
        if not self.id:
            self.id = uuid.uuid4()
        super().save(*args, **kwargs)
    
    def get_all_bookings_for_user(user: User):
        return BookingTransaction.objects.filter(user=user)