require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("@fortawesome/fontawesome-free/js/all.js");

$(() => {
    $('[data-toggle="tooltip"]').tooltip()

    $('.clickable').click(function(e) {
        var tr = $(e.target).closest(".clickable");
        window.location = tr.data("href");
    })
});
