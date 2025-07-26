document.querySelectorAll('.truncate-cell').forEach(cell => {
  const tooltip = document.createElement('div');
  tooltip.className = 'custom-tooltip';

  // Gamitin .innerHTML instead of .innerText para ma-render ang HTML tags
  let fullText = cell.dataset.fulltext || cell.innerText;

  // Palitan lang ang [MM-DD-YYYY HH:MM:SS AM/PM] sa simula ng line (optional: multiple)
  fullText = fullText.replace(
    /\[(\d{2}-\d{2}-\d{4} \d{2}:\d{2}:\d{2} [AP]M)\]/g,
    '<span style="color:red">[$1]</span>'
  );

  tooltip.innerHTML = fullText;

  // Inline styles para wala nang kailangan CSS file
  Object.assign(tooltip.style, {
    display: 'none',
    position: 'fixed',
    background: '#2b2b2b',
    color: '#fff',
    padding: '10px',
    fontSize: '0.875rem',
    borderRadius: '6px',
    boxShadow: '0 4px 14px rgba(0, 0, 0, 0.3)',
    zIndex: '9999',
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-word',
    maxWidth: '300px',
    maxHeight: '250px',
    overflowY: 'auto'
  });

  document.body.appendChild(tooltip);

  let hideTimeout;

  const showTooltip = () => {
    tooltip.style.display = 'block';

    tooltip.style.visibility = 'hidden';
    tooltip.style.left = '0px';
    tooltip.style.top = '0px';

    const rect = cell.getBoundingClientRect();
    const tooltipRect = tooltip.getBoundingClientRect();

    let top = rect.top + window.scrollY - tooltipRect.height - 8;
    if (top < window.scrollY) {
      top = rect.bottom + window.scrollY + 8;
    }

    let left = rect.left + window.scrollX;
    if (left + tooltipRect.width > window.innerWidth) {
      left = window.innerWidth - tooltipRect.width - 10;
    }

    tooltip.style.left = left + 'px';
    tooltip.style.top = top + 'px';
    tooltip.style.visibility = 'visible';
  };

  const hideTooltip = () => {
    hideTimeout = setTimeout(() => {
      tooltip.style.display = 'none';
    }, 150);
  };

  cell.addEventListener('mouseenter', showTooltip);
  cell.addEventListener('mouseleave', hideTooltip);

  tooltip.addEventListener('mouseenter', () => {
    clearTimeout(hideTimeout);
  });

  tooltip.addEventListener('mouseleave', hideTooltip);
});
