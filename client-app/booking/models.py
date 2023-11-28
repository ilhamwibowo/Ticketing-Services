import uuid
from django.contrib.auth.models import User
from django.db import models

class BookingTransaction(models.Model):
    STATUS_CHOICES = (
        ('SUCCESS', 'Success'),
        ('FAILED', 'Failed'),
        ('PENDING', 'Pending'),
    )
    id = models.UUIDField(primary_key=True, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    event_id = models.CharField(max_length=100)
    seats = models.JSONField(default=list)  # List of booked seats
    status = models.CharField(max_length=10, choices=STATUS_CHOICES)

    def save(self, *args, **kwargs):
        # Generate UUID if not provided
        if not self.id:
            self.id = uuid.uuid4()
        super().save(*args, **kwargs)

class Invoice(models.Model):
    id = models.CharField(primary_key=True, max_length=100)  # ID as a string
    transaction = models.ForeignKey(BookingTransaction, on_delete=models.CASCADE)
    invoice = models.FileField()  # Store the PDF file

    class Meta:
        unique_together = ("transaction",)

    def __str__(self):
        return f"Invoice {self.id} for Transaction {self.transaction.id}"
    
    def save(self, *args, **kwargs):
        # Generate UUID if not provided
        if not self.id:
            self.id = str(uuid.uuid4())
        super().save(*args, **kwargs)