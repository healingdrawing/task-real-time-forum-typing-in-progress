import { VIEWS } from './view.mjs';

export class WSClient {
	constructor(forum) {
		this.forum = forum;
		this.socket = new WebSocket('ws://localhost:8080/ws');
		this.socket.onmessage = (event) => {
			const data = JSON.parse(event.data);
			switch (data.type) {
				case 'post':
					if (this.forum.state.page === VIEWS.POSTS)
						this.forum.view.posts.insertBefore(
							this.forum.view.generatePost(data.post),
							this.forum.view.posts.firstChild
						);
					break;
				case 'comment':
					if (this.forum.state.page === VIEWS.POST && data.postID === this.forum.view.viewPost.postID)
						this.forum.view.comments.insertBefore(
							this.forum.view.generateComment(data.comment),
							this.forum.view.comments.firstChild
						);
					break;
				case 'status':
					if (this.forum.view.users[data.username]) {
						this.forum.view.users[data.username].data.online = data.online;
						this.forum.view.users[data.username].element = this.forum.view.generateUser(
							this.forum.view.users[data.username].data
						);
						this.forum.view.users[data.username].data.typing = false;
						this.forum.view.updateTyping();
					} else {
						this.forum.view.users[data.username] = {
							data,
							element: this.forum.view.generateUser(data),
							typing: false,
						};
					}
					this.forum.view.updateUserList();

					break;
				case 'message':
					const target = this.forum.state.auth.username == data.message.usernameFrom ? data.message.usernameTo : data.message.usernameFrom;
					this.forum.view.users[target].data.time = data.message.time;
					const sender = data.message.usernameTo;
					if (this.forum.view.users[sender]) this.forum.view.users[sender].data.time = data.message.time;
					this.forum.view.users[target].element = this.forum.view.generateUser(
						this.forum.view.users[target].data
					);
					this.forum.view.updateUserList();

					if (data.message.usernameFrom != this.forum.state.auth.username && (this.forum.state.page != VIEWS.CHAT ||
						(!this.forum.view.unread.includes(target) && !(!this.forum.view.viewMessages.hidden && this.forum.view.viewMessages.username == target)) ||
						this.forum.view.messages.scrollTop + this.forum.view.messages.clientHeight != this.forum.view.messages.scrollHeight)) {
						this.forum.view.unread.push(target);
						this.forum.view.updateNotifications();
						this.forum.view.updateUserList();
					}
					if (
						!this.forum.view.viewMessages.hidden &&
						target == this.forum.view.viewMessages.username
					) {
						const iniUsrmsg = this.forum.view.messages.innerText.split("\n")[0]
						const iniShowmsg = "There are no messages between you and "
						if (this.forum.view.messages.innerText.includes(iniShowmsg)) {
							this.forum.view.messages.replaceChildren(this.forum.view.generateMessage(data.message));
						}
						else if (
							this.forum.view.messages.scrollTop + this.forum.view.messages.clientHeight == this.forum.view.messages.scrollHeight
						) {
							this.forum.view.messages.appendChild(this.forum.view.generateMessage(data.message));
							this.forum.view.messages.scrollTop = this.forum.view.messages.scrollHeight;
						} else this.forum.view.messages.appendChild(this.forum.view.generateMessage(data.message));
						this.forum.view.messageBank.push(this.forum.view.generateMessage(data.message));
					}
					this.updateScroll();
					break;
				case 'typing':
					this.forum.view.users[data.usernameFrom].data.typing = data.typing
					this.forum.view.updateTyping();
					break
			}
		};
	}

	updateScroll() {
		var element = document.getElementById("messages");
		element.scrollTop = element.scrollHeight;
	}
	

	login() {
		this.socket.send(JSON.stringify({ type: 'login', data: { username: this.forum.state.auth.username } }));
	}

	logout() {
		this.socket.send(JSON.stringify({ type: 'logout', data: { username: this.forum.state.auth.username } }));
	}
}

