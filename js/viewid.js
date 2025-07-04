  $(document).ready(function() {
    $('.editBtn').on('click', function() {
        const btn = $(this);
        $('#modalID').val(btn.data('id'));
        $('#modalLink').val(btn.data('link'));
     
    });
  });


  $('.editBtn').on('click', function () {
    const btn = $(this);
    const id = btn.data('id');

    $('#modalIDText').text(id);   // for display
    $('#modalID').val(id);        // for form POST
});
