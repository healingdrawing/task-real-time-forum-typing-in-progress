export const throttle = (func, wait, options = { leading: false, trailing: false }) => {
	let timeout = null;
	let stash = null;

	const delay = () => {
		timeout = setTimeout(lookAtStash, wait);
	};

	const lookAtStash = () => {
		timeout = null;
		if (stash !== null) {
			func(...stash);
			stash = null;
			delay();
		}
	};

	return function (...args) {
		if (timeout !== null) {
			if (options.trailing) {
				stash = args;
			}
			return;
		}

		if (options.leading) {
			func(...args);
			delay();
			return;
		}

		if (options.trailing) {
			stash = args;
			delay();
		}
	};
};

export function showAlert(message) {
	// Get the modal box and the close button
	const modal = document.querySelector('.modal');
	const closeBtn = document.querySelector('.modal-close');

	// Set the message content of the modal box
	const modalContent = modal.querySelector('.modal-content');
	modalContent.textContent = message;

	// Show the modal box and the close button
	modal.style.display = 'block';
	modal.style.opacity = 1;
	modal.style.height = '100%';

	if (closeBtn) closeBtn.style.display = 'block';

	const handleClick = (event) => {
		// Check if the click event target is outside the modal
		if (!modalContent.contains(event.target)) {
			// Hide the modal box and the close button
			modal.style.opacity = 0;
			modal.style.height = '0%';
			setTimeout(function () {
				modal.style.display = 'none';
				if (closeBtn) closeBtn.style.display = 'none';
			}, 10);

			document.body.removeEventListener('click', handleClick);
		}
	};

	document.body.addEventListener('click', handleClick);

	// Hide the modal box and the close button after 3.5 seconds
	setTimeout(function () {
		modal.style.opacity = 0;
		modal.style.height = '0%';
		setTimeout(function () {
			modal.style.display = 'none';
			if (closeBtn) closeBtn.style.display = 'none';
		}, 300);
		document.body.removeEventListener('click', handleClick);
	}, 3500);
}