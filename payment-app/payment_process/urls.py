from django.urls import path
from .views import process_payment, payment_form

urlpatterns = [
    path('process-payment/', process_payment, name='process_payment'),
    path('payment-form/', payment_form, name='payment_form'),
]