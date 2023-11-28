from django.db import models

# Create your models here.
class Invoice(models.Model):
    invoice_number = models.CharField(max_length=255)
    amount = models.DecimalField(max_digits=10, decimal_places=2)
    success = models.BooleanField(default=False)