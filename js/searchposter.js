$(document).ready(function () {
        $('#posterTable').DataTable({
            "order": [], // Default: no initial sort
            "language": {
                "search": "🔍 Search:"
            }
        });

        // Handle modal population
        $('.editBtn').click(function () {
            var id = $(this).data('id');
            var link = $(this).data('link');
            $('#modalID').val(id);
            $('#modalIDText').text(id);
            $('#modalLink').val(link);
        });
    });