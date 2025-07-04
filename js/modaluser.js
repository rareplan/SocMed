$('#editModal').on('show.bs.modal', function (event) {
    const button = $(event.relatedTarget);
    const modal = $(this);

    modal.find('#modalID').val(button.data('id'));
    modal.find('#modalUsername').val(button.data('username'));
    modal.find('#modalPassword').val(button.data('password'));
    modal.find('#modalRole').val(button.data('role'));
  });