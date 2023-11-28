from django.urls import path, include
from django.contrib.auth.decorators import login_required

from .views import list_bookings, refresh_booking_status, BookView, get_available_events, get_chairs_status, list_of_bookings

urlpatterns = [
    path("", list_of_bookings, name='bookings'),
    path("bookings/", list_bookings, name='list_bookings'),
    path("bookings/<int:booking_id>/refresh/", refresh_booking_status, name='refresh_booking_status'),

    path('book/', login_required(BookView.as_view()), name='book'),

    path('events/', get_available_events, name='get_available_events'),
    path('events/<str:event_id>/chairs/', get_chairs_status, name='get_chairs_status'),
]