import { showAlert } from "./utilities.mjs";

window.alert = function (message) {
	showAlert(message);
};

export const checkSignup = () => {
	const data = new FormData(document.getElementById('view-signup')); // get the data from the html form
	const requirements = {
		Nickname: { // the name of the field will be used to display the error message
			field: 'signup-user', // the id of the field in the html form, and in data
			check: (value) => value.match(/^\w{2,15}$/i), // the function that checks the value of the field
			requirement: `2-15 english letters/numbers`, // the requirement that will be displayed in the error message
		},
		Age: {
			field: 'signup-age',
			check: (value) => value !== '' && Number.parseInt(value, 10) > 0 && Number.parseInt(value, 10) < 1000,
			requirement: 'We expect you to be between 1 to 999 years old',
		},
		'First name': {
			field: 'signup-firstname',
			check: (value) => value.match(/^[\w-]{1,32}$/i),
			requirement: '1-32 english letters/numbers',
		},
		'Last name': {
			field: 'signup-lastname',
			check: (value) => value.match(/^[\w-]{1,32}$/i),
			requirement: '1-32 english letters/numbers',
		},
		'E-mail': {
			field: 'signup-email',
			check: (value) => value.match(/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,4})+$/i),
			requirement: 'valid e-mail',
		},
		Password: {
			field: 'signup-password',
			check: (value) => value.match(/^\w{6,15}$/i),
			requirement: '6-15 english letters/numbers',
		},
	};
	for (const [fieldName, checkObj] of Object.entries(requirements))
		if (!checkObj.check(data.get(checkObj.field))) { // check the value of the field
			alert(`${fieldName} requirement: ${checkObj.requirement}`);
			return false;
		}
	return true;
};

export const checkLogin = () => {
	const data = new FormData(document.getElementById('view-login'));
	const requirements = {
		Nickname: {
			field: 'login-user',
			check: (value) => value.match(/(^\w{2,15}$)|(^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,4})+$)/i),
			requirement: '2-15 english letters/numbers',
		},
		Password: {
			field: 'login-password',
			check: (value) => value.match(/^\w{6,15}$/i),
			requirement: '6-15 english letters/numbers',
		},
	};
	for (const [fieldName, checkObj] of Object.entries(requirements))
		if (!checkObj.check(data.get(checkObj.field))) {
			alert(`${fieldName} is not meeting requirements:
		required: ${checkObj.requirement}`);
			return false;
		}
	return true;
};

export const checkPost = () => {
	const data = new FormData(document.getElementById('view-newpost'));
	const requirements = {
		Title: {
			field: 'newpost-title',
			check: (value) => !!value,
			requirement: 'cannot be empty',
		},
		Text: {
			field: 'newpost-text',
			check: (value) => !!value,
			requirement: 'cannot be empty',
		},
	};
	for (const [fieldName, checkObj] of Object.entries(requirements))
		if (!checkObj.check(data.get(checkObj.field))) {
			alert(`${fieldName} requirement: ${checkObj.requirement}`);
			return false;
		}
	return true;
};

export const checkComment = () => {
	const data = new FormData(document.getElementById('comment-form'));
	const requirements = {
		Text: {
			field: 'comment-text',
			check: (value) => !!value,
			requirement: 'cannot be empty',
		},
	};
	for (const [fieldName, checkObj] of Object.entries(requirements))
		if (!checkObj.check(data.get(checkObj.field))) {
			alert(`${fieldName} requirement: ${checkObj.requirement}`);
			return false;
		}
	return true;
};

export const checkMessage = () => {
	const data = new FormData(document.getElementById('message-form'));
	const requirements = {
		Message: {
			field: 'message-text',
			check: (value) => !!value,
			requirement: 'cannot be empty',
		},
	};
	for (const [fieldName, checkObj] of Object.entries(requirements))
		if (!checkObj.check(data.get(checkObj.field))) {
			alert(`${fieldName} requirement: ${checkObj.requirement}`);
			return false;
		}
	return true;
};
