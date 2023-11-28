from django.shortcuts import render
from django.contrib.auth.forms import UserCreationForm
from django.contrib.auth import login
from django.http import HttpResponseRedirect

def register_user(request):
    if request.method == 'POST':
        form = UserCreationForm(request.POST)
        if form.is_valid():
            user = form.save()
            # Log in the user after registration
            login(request, user)
            return HttpResponseRedirect('/')  # Redirect to profile page after successful registration
    else:
        form = UserCreationForm()
    
    return render(request, 'registration/register.html', {'form': form})
