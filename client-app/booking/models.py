import uuid
from django.contrib.auth.models import User
from django.db import models

class BookingTransaction(models.Model):
    STATUS_CHOICES = (
        ('SUCCESS', 'Success'),
        ('FAILED', 'Failed'),
        ('PENDING', 'Pending'),
    )
    uuid = models.UUIDField(primary_key=True, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    event_id = models.CharField(max_length=100)
    seats = models.JSONField(default=list)  # List of booked seats
    status = models.CharField(max_length=10, choices=STATUS_CHOICES)

    def save(self, *args, **kwargs):
        # Generate UUID if not provided
        if not self.uuid:
            self.uuid = uuid.uuid4()
        super().save(*args, **kwargs)