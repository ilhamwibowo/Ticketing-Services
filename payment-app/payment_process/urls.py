from django.urls import path
from .views import send_invoice_id, process_payment, payment_form

urlpatterns = [
    path('process-payment/', process_payment, name='process_payment'),
]