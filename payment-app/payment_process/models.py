from django.db import models

# Create your models here.
class Invoice(models.Model):
    invoice_id = models.CharField(max_length=255)
    status = models.BooleanField(default=False)