from django.urls import path
from .views import process_payment

urlpatterns = [
    path('process-payment/', process_payment, name='process_payment'),
]