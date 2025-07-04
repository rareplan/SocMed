 function goBack() {
    location.href = '/allposterchecker';
    }

    const colors = ['#ff4081', '#4caf50', '#ffc107', '#03a9f4', '#ff5722', '#e91e63', '#00e676'];
    const confetti = document.getElementById('confetti');

    for (let i = 0; i < 100; i++) {
      const square = document.createElement('div');
      square.classList.add('square');
      square.style.left = Math.random() * 100 + 'vw';
      square.style.top = Math.random() * -100 + 'px';
      square.style.backgroundColor = colors[Math.floor(Math.random() * colors.length)];
      square.style.animationDuration = (5 + Math.random() * 5) + 's';
      confetti.appendChild(square);
    }