const $ = document.querySelector.bind(document);
const $$ = document.querySelectorAll.bind(document);

document.addEventListener('DOMContentLoaded', () => {
  const imgInput = $('#imgInput');
  const mode = $('#modeOption');
  const shape = $('#shapeOption');
  const output = $('#typeOption');
  const btn = $('button');
  const imgResult = $('#imgResult');

  btn.addEventListener('click', () => {
    const url = `${window.location.origin}/primitive.go?img=${
      imgInput.value
    }&mode=${mode[mode.selectedIndex].value}&shape=${
      shape[shape.selectedIndex].value
    }&output=${output[output.selectedIndex].value}`;

    disableForm(true);

    imgResult.setAttribute('src', '');
    imgResult.setAttribute('src', url);

    imgResult.addEventListener('load', () => {
      disableForm(false);
    });
  });

  function disableForm(disable) {
    const inputs = $$('.js-input');
    const btn = $('.js-button');
    const imgResult = $('#imgResult');
    const loading = $('.js-loading');

    if (disable) {
      btn.setAttribute('disabled', true);
      imgResult.classList.add('hide');
      loading.classList.remove('hide');
      return inputs.forEach(input => {
        input.setAttribute('disabled', 'true');
      });
    }

    loading.classList.add('hide');
    btn.removeAttribute('disabled');
    imgResult.classList.remove('hide');
    inputs.forEach(input => {
      input.removeAttribute('disabled');
    });
  }
});
