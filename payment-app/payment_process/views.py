import json
import random
import requests
import asyncio
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.urls import reverse
from .models import Invoice

@csrf_exempt
async def process_payment(request):
    invoice_id = generate_unique_invoice_id()

    # Simulasi kegagalan 10%
    status = await random.choices([True, False], weights=[90, 10])[0]

    invoice = Invoice.objects.create(
        invoice_id=invoice_id,
        status=status,
    )

    # Kirim webhook ke Ticket App
    webhook_url = "http://localhost:3000/webhook/payment"  # Ganti dengan konfigurasi yang sesuai
    payment_url = request.build_absolute_uri(reverse('process_payment'))

    payload = {
        'invoice_id': invoice.invoice_id,
        'status': invoice.status,
    }

    loop = asyncio.get_event_loop()
    await loop.run_in_executor(None, lambda: requests.post(webhook_url, data=json.dumps(payload), headers={'Content-Type': 'application/json'}))
    
    return JsonResponse({
        'invoice_id': invoice_id, 
        'payment_url': payment_url
    })

async def generate_unique_invoice_id():
    # Fungsi ini harus memastikan invoice_id adalah unik dan belum ada di database
    while True:
        invoice_id = generate_random_invoice_id()
        if not Invoice.objects.filter(invoice_id=invoice_id).exists():
            return invoice_id

def generate_random_invoice_id():
    # Implementasi sederhana untuk menghasilkan invoice_id yang unik
    return ''.join(random.choices('ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', k=8))

# def payment_form(request):
#     invoice_id = request.GET.get('invoice_id')
#     return render(request, 'payment_form.html', {'invoice_id': invoice_id})