const $ = document.querySelector.bind(document);
const $$ = document.querySelectorAll.bind(document);

document.addEventListener('DOMContentLoaded', () => {
  const imgInput = $('#imgInput');
  const mode = $('#modeOption');
  const shape = $('#shapeOption');
  const btn = $('button');
  const imgResult = $('#imgResult');

  btn.addEventListener('click', () => {
    if (imgInput.value === 'https://sophearak.me/static/profile.jpg') {
      const params = `https://image-to-primitive-gd1lniddr.now.sh/?img=${
        imgInput.value
      }&mode=${mode[mode.selectedIndex].value}&shape=${
        shape[shape.selectedIndex].value
      }`;

      disableForm(true);

      imgResult.setAttribute('src', '');
      imgResult.setAttribute('src', params);

      imgResult.addEventListener('load', () => {
        disableForm(false);
      });
    }
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

  // imgInput.addEventListener('blur', event => {
  //   if (event.target.value !== 'https://sophearak.me/static/profile.jpg') {
  //     getPrimitive();
  //   }
  // });

  // function getPrimitive() {
  //   imgResult.setAttribute('src', '');
  //   const params = `https://image-to-primitive-gd1lniddr.now.sh/?img=${
  //     imgInput.value
  //   }&mode=${mode[mode.selectedIndex].value}&shape=${
  //     shape[shape.selectedIndex].value
  //   }`;
  //   imgResult.setAttribute('src', params);
  // }
});
