/**
 * A component for displaying a temporary toast message
 *
 * @example
 * ```html
 * <toast-message data-message="HELLO"></toast-message>
 * ```
 */
class ToastMessage extends HTMLElement {
  /** @type {HTMLElement | null}*/
  #wrapper = null;
  /** @type {HTMLElement | null}*/
  #button = null;
  /** @type {number | null}*/
  #closeTimeoutID;
  /** 
    * Used to be able to pass #close as a callback
    * @type {Function}
    */
  #bindClose = this.#close.bind(this);

  constructor() {
    super();
  }

  connectedCallback() {
    const shadow = this.attachShadow({ mode: "open" });
    this.#wrapper = document.createElement("figure");
    this.#wrapper.setAttribute("class", "toast-message closed");

    const content = document.createElement("figcaption");
    content.setAttribute("class", "toast-content");
    content.textContent = this.getAttribute("data-message");

    this.#button = document.createElement("button");
    this.#button.innerHTML = "&#10005;";

    const style = document.createElement("style");

    style.textContent = ToastMessage.STYLES;

    this.#wrapper.appendChild(content);
    this.#wrapper.appendChild(this.#button);
    shadow.appendChild(style);
    shadow.appendChild(this.#wrapper);

    this.#button.addEventListener("click", this.#bindClose);

    setTimeout(() => {
      this.#open();
      //this.#scheduleToClose();
    }, 100);

    /*
    <figure id="toast-message" class="toast-message">
      <figcaption id="message-content" class="toast-contnent">{{.}}</figcaption>
    </figure>
      */
  }

  #open() {
    if (!this.#wrapper) return;
    this.#wrapper.setAttribute("class", "toast-message open");
  }

  #close() {
    if (!this.#wrapper) return;
    this.#wrapper.setAttribute("class", "toast-message closed");
    this.#button.removeEventListener("click", this.#bindClose);
    clearTimeout(this.#closeTimeoutID);
  }

  #scheduleToClose() {
    if (!this.#wrapper) return;
    this.#closeTimeoutID = setTimeout(this.#bindClose, 3000);
  }

  static STYLES = `
    * {
      margin: 0;
      padding: 0;
    }

    .toast-message {
      display: flex;
      gap: 1rem;
      align-items: center;
      justify-content: space-between;
      position: absolute;
      bottom: 2rem;
      min-width: 15rem;
      padding: 2rem;
      border-radius: 5px;
      background-color: var(--murrey-100);
      color: var(--ghost-white-100);
      transition: transform .1s;

      &.closed {
        left: 0;
        transform: translateX(-100%);
      }

      &.open {
        left: 2rem;
        transform: translateX(0);
      }

      button {
        border: none;
        cursor: pointer;
        background-color: transparent;
        color: var(--ghost-white-100);
      }
    }
  `;
}

customElements.define("toast-message", ToastMessage);
