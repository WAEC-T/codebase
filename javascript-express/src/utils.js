const formatMessages = (messages, isApi = false) => {
    return messages.map((message) => {
        const updatedMessage = { ...message };

        Object.entries(updatedMessage).forEach(([key, value]) => {
            if (key.startsWith('User.')) {
                const newKey = key.replace('User.', '');
                updatedMessage[newKey] = value;
                delete updatedMessage[key];
            }
            if (isApi && key === 'text') {
                updatedMessage['content'] = value;
                delete updatedMessage[key];
            }
        });

        return updatedMessage;
    });
};

const validateRegisterFields = async (
    username,
    email,
    password,
    password2,
    isApi = false
) => {
    let errorMessage = null;

    if (!username) {
        errorMessage = 'You have to enter a username.';
    } else if (!email || !email.includes('@')) {
        errorMessage = 'You have to enter a valid email address.';
    } else if (!password) {
        errorMessage = 'You have to enter a password.';
    } else if (password !== password2 && !isApi) {
        errorMessage = 'The two passwords do not match.';
    }
    return errorMessage;
};

module.exports = { formatMessages, validateRegisterFields };
