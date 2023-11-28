from django.urls import path, include
from django.contrib.auth.decorators import login_required

from .views import (
    list_bookings, 
    refresh_booking_status, 
    BookView, 
    get_available_events, 
    get_chairs_status, 
    list_of_bookings,
    create_invoice,
    get_invoice
)

urlpatterns = [
    path("", list_of_bookings, name='bookings'),
    path("api/bookings/", list_bookings, name='list_bookings'),
    path("api/bookings/<str:booking_id>/refresh/", refresh_booking_status, name='refresh_booking_status'),

    path('book/', login_required(BookView.as_view()), name='book'),

    path('api/events/', get_available_events, name='get_available_events'),
    path('api/events/<str:event_id>/chairs/', get_chairs_status, name='get_chairs_status'),

    path('api/invoices/create/', create_invoice, name='create_invoice'),
    path('api/invoices/<str:invoice_id>/', get_invoice, name='get_invoice'),
]