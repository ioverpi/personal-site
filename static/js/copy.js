// Run immediately since this script is loaded on the page that needs it
(function() {
  const copyBtn = document.getElementById('copy-invite-btn');
  if (copyBtn) {
    copyBtn.addEventListener('click', function() {
      const input = document.getElementById('invite-url');
      input.select();
      navigator.clipboard.writeText(input.value);
    });
  }
})();
