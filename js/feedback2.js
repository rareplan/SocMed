document.querySelectorAll('.truncate-cell').forEach(cell => {
  const tooltip = document.createElement('div');
  tooltip.className = 'custom-tooltip';

  let fullText = cell.dataset.fulltext || cell.innerText;

  // Kulayan ang timestamp
  fullText = fullText.replace(
    /\[(\d{2}-\d{2}-\d{4} \d{2}:\d{2}:\d{2} [AP]M)\]/g,
    '<span style="color:red">[$1]</span>'
  );

  tooltip.innerHTML = fullText;

  // Attach sa pinakamalapit na .table-responsive
  const container = cell.closest('.table-responsive') || document.body;
  container.appendChild(tooltip);

  // Gamit absolute positioning sa loob ng scrollable container
  Object.assign(tooltip.style, {
    display: 'none',
    position: 'absolute',
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

  let hideTimeout;

  const showTooltip = () => {
    tooltip.style.display = 'block';

    tooltip.style.visibility = 'hidden';
    tooltip.style.left = '0px';
    tooltip.style.top = '0px';

    const rect = cell.getBoundingClientRect();
    const containerRect = container.getBoundingClientRect();

    const tooltipRect = tooltip.getBoundingClientRect();

    let top = rect.top - containerRect.top - tooltipRect.height - 8 + container.scrollTop;
    if (top < 0) {
      top = rect.bottom - containerRect.top + 8 + container.scrollTop;
    }

    let left = rect.left - containerRect.left + container.scrollLeft;
    if (left + tooltipRect.width > container.clientWidth) {
      left = container.clientWidth - tooltipRect.width - 10 + container.scrollLeft;
    }

    tooltip.style.top = `${top}px`;
    tooltip.style.left = `${left}px`;
    tooltip.style.visibility = 'visible';
  };

  const hideTooltip = () => {
    hideTimeout = setTimeout(() => {
      tooltip.style.display = 'none';
    }, 150);
  };

  cell.addEventListener('mouseenter', showTooltip);
  cell.addEventListener('mouseleave', hideTooltip);
  tooltip.addEventListener('mouseenter', () => clearTimeout(hideTimeout));
  tooltip.addEventListener('mouseleave', hideTooltip);
});
