/**
 * 指定したメッセージをトースト通知として表示する
 * @param {string} message 
 * @param {string} type 'success', 'danger', 'warning', 'info' など
 */
function showToast(message, type = 'success') {
    let container = document.getElementById('toastContainer');
    
    // コンテナがなければ作成してbodyに追加
    if (!container) {
        container = document.createElement('div');
        container.id = 'toastContainer';
        container.className = 'toast-container position-fixed bottom-0 end-0 p-3';
        document.body.appendChild(container);
    }

    const toastEl = document.createElement('div');
    toastEl.className = `toast align-items-center border-0 text-bg-${type}`;
    toastEl.setAttribute('role', 'alert');
    toastEl.setAttribute('aria-live', 'assertive');
    toastEl.setAttribute('aria-atomic', 'true');

    toastEl.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
    `;

    container.appendChild(toastEl);
    const toast = new bootstrap.Toast(toastEl);
    toast.show();

    toastEl.addEventListener('hidden.bs.toast', () => {
        toastEl.remove();
    });
}

/**
 * クリップボードにテキストをコピーし、ボタンにフィードバックを表示する
 * @param {string|HTMLInputElement} input コピー対象のIDまたは要素
 * @param {HTMLElement} btn フィードバックを表示するボタン要素
 */
function copyToClipboard(input, btn) {
    const inputEl = typeof input === 'string' ? document.getElementById(input) : input;
    const textToCopy = inputEl.value || inputEl.textContent;
    const originalContent = btn.innerHTML;

    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(textToCopy).then(() => {
            showCopyFeedback(btn, originalContent);
        }).catch(err => {
            console.error("Clipboard API failed:", err);
            fallbackCopy(textToCopy, btn, originalContent);
        });
    } else {
        fallbackCopy(textToCopy, btn, originalContent);
    }
}

function fallbackCopy(text, btn, originalContent) {
    const textArea = document.createElement("textarea");
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.select();
    try {
        document.execCommand('copy');
        showCopyFeedback(btn, originalContent);
    } catch (err) {
        console.error("Fallback copy failed:", err);
        showToast("コピーに失敗しました。", "danger");
    }
    document.body.removeChild(textArea);
}

function showCopyFeedback(btn, originalContent) {
    btn.innerHTML = '<i class="bi bi-check-lg"></i> Copied!';
    const originalClass = btn.className;
    
    // 一時的にボタンの色を変える（Bootstrapのクラスを想定）
    if (btn.classList.contains('btn-outline-secondary')) {
        btn.classList.replace('btn-outline-secondary', 'btn-success');
    } else if (btn.classList.contains('btn-primary')) {
        btn.classList.replace('btn-primary', 'btn-success');
    }

    showToast("クリップボードにコピーしました", "success");

    setTimeout(() => {
        btn.innerHTML = originalContent;
        btn.className = originalClass;
    }, 2000);
}
