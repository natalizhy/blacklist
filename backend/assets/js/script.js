$(".toggle-menu-btn").on('click', function () {
    $(`.toggle-menu-btn .show`).removeClass('show');
    $(this).toggleClass('show');
});

document.querySelector("input").addEventListener("change", function () {
    if (this.files[0]) {
        var fr = new FileReader();

        fr.addEventListener("load", function () {
            document.querySelector(".update").classList.add('hoverable');
            document.querySelector(".img-bg").style.backgroundImage = `url("${fr.result}")`;
        }, false);

        fr.readAsDataURL(this.files[0]);
    }
});

var modal = document.getElementById('user');

var recaptchaCallback = function () {
    $('#bt').css('background-color', '#2196F3');
    $('#bt').removeAttr('disabled');
};