(function () {
  // hero: the group tabs (you write / generated / run) and, inside each group, the file list
  function tablist(list) {
    var vertical = list.getAttribute('aria-orientation') === 'vertical';
    var prev = vertical ? 'ArrowUp' : 'ArrowLeft';
    var next = vertical ? 'ArrowDown' : 'ArrowRight';
    var tabs = Array.prototype.filter.call(list.children, function (el) { return el.getAttribute('role') === 'tab'; });
    function select(tab) {
      tabs.forEach(function (t) {
        var on = t === tab;
        t.setAttribute('aria-selected', on ? 'true' : 'false');
        t.tabIndex = on ? 0 : -1;
        document.getElementById(t.getAttribute('aria-controls')).hidden = !on;
      });
    }
    tabs.forEach(function (t, i) {
      t.addEventListener('click', function () { select(t); });
      t.addEventListener('keydown', function (e) {
        var n = e.key === next ? tabs[(i + 1) % tabs.length]
          : e.key === prev ? tabs[(i + tabs.length - 1) % tabs.length]
          : e.key === 'Home' ? tabs[0]
          : e.key === 'End' ? tabs[tabs.length - 1] : null;
        if (n) { select(n); n.focus(); e.preventDefault(); }
      });
    });
  }
  document.querySelectorAll('#hero-win [role="tablist"]').forEach(tablist);

  // copy buttons
  document.querySelectorAll('.copy').forEach(function (b) {
    b.addEventListener('click', function () {
      var el = document.getElementById(b.getAttribute('data-copy'));
      var text = el.textContent;
      function done() { b.textContent = 'copied'; b.classList.add('ok'); setTimeout(function () { b.textContent = 'copy'; b.classList.remove('ok'); }, 1600); }
      function fallback() { var r = document.createRange(); r.selectNodeContents(el); var s = getSelection(); s.removeAllRanges(); s.addRange(r); b.textContent = 'selected'; }
      try {
        navigator.clipboard.writeText(text).then(done, fallback);
      } catch (e) { fallback(); }
    });
  });
})();
