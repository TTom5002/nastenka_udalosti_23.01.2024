// function Goback() {
//     document.getElementById("go-back").addEventListener("click", () => {
//         const lastVisitedPage = sessionStorage.getItem("lastVisitedPage");
//         if (lastVisitedPage) {
//             window.location.href = lastVisitedPage;
//             sessionStorage.removeItem("lastVisitedPage"); // Odstranění hodnoty po použití
//         } else {
//             history.back(); // Jako záložní plán, pokud URL nebyla uložena
//         }
//     });
// }

// function LoadMore() {
//     let offset = 0; // Startovací bod pro načítání příspěvků
//     const limit = 2; // Kolik příspěvků chcete načíst najednou

//     document.getElementById('load-more').addEventListener('click', function () {
//         offset += limit; // Zvyšte offset pro načtení dalších příspěvků
//         fetch(`/?limit=${limit}&offset=${offset}`)
//             .then(response => response.json())
//             .then(data => {
//                 // Předpokládáme, že 'data' je pole objektů 'event'
//                 const container = document.getElementById('event-container');
//                 data.forEach(event => {
//                     const eventElement = document.createElement('div');
//                     eventElement.className = 'event';
//                     eventElement.innerHTML = `<h2>${event.Header}</h2><p>${event.Body}</p>`;
//                     // Přidání dalších detailů eventu podle potřeby
//                     container.appendChild(eventElement);
//                 });
//             })
//             .catch(error => console.error('Chyba při načítání více příspěvků:', error));
//     });
// }

function Prompt() {
    let toast = function (c) {
        const {
            msg = '',
            icon = 'success',
            position = 'top-end',

        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

// document.getElementById("go-back").addEventListener("click", () => {
//     history.back();
// });

