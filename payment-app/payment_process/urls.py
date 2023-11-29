from django.urls import path
from .views import process_payment, hello_world

urlpatterns = [
    path('process-payment/', process_payment, name='process_payment'),
    path('hello-world/', hello_world, name='hello_world'),
]