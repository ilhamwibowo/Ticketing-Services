from django.urls import path
from .views import send_invoice_id, process_payment, payment_form

urlpatterns = [
    path('send-invoice-id/', send_invoice_id, name='send_invoice_id'),
    path('process-payment/', process_payment, name='process_payment'),
    path('payment-form/', payment_form, name='payment_form'),
]