 function submitWithAnimation() {
      const box = document.getElementById('confirmBox');
      const spinner = document.getElementById('spinner');
      const form = document.getElementById('logoutForm');

      box.classList.add('fade-out');
      spinner.style.display = 'block';

      setTimeout(() => {
        form.submit(); // Actually submit the logout form
      }, 1500);
    }

    function cancelLogout() {
      window.history.back(); // Cancel logout and go back
    }