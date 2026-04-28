document.body.addEventListener('htmx:responseError', function () {
  console.warn('Error en solicitud HTMX');
});
