import json
import random
import requests
import asyncio
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.urls import reverse
from .models import Invoice
from django.shortcuts import render

@csrf_exempt
def send_invoice_id(request):
    while True:
        invoice_id = ''.join(random.choices('ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', k=8))
        if not Invoice.objects.filter(invoice_id=invoice_id).exists():
            return invoice_id, request.build_absolute_uri(reverse('payment_form')) + f'?invoice_id={invoice_id}'

@csrf_exempt
async def process_payment(request):
    invoice_id = request.POST.get('invoice_id')

    # Simulasi kegagalan 10%
    status = await random.choices([True, False], weights=[90, 10])[0]

    invoice = Invoice.objects.create(
        invoice_id=invoice_id,
        status=status,
    )

    # Kirim webhook ke Ticket App
    webhook_url = "WEBHOOK_URL"  # Ganti dengan konfigurasi yang sesuai
    payload = {
        'invoice_id': invoice.invoice_id,
        'status': invoice.status,
    }

    loop = asyncio.get_event_loop()
    await loop.run_in_executor(None, lambda: requests.post(webhook_url, data=json.dumps(payload), headers={'Content-Type': 'application/json'}))
    
    return JsonResponse({'message': 'Payment process done!'})

def payment_form(request):
    invoice_id = request.GET.get('invoice_id')
    return render(request, 'payment_form.html', {'invoice_id': invoice_id})