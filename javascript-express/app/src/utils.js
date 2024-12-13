const formatMessages = (messages) => {
    const formattedMessages = messages.map((message) => {
      return {
        user_id: message.author_id,
        username: message.username,
        email: message.email,
        message_id: message.message_id,
        text: message.text,
        pub_date: new Date(message.pub_date).toISOString().replace(/T/, ' @ ').slice(0, 16),
        flagged: message.flagged,
      };
    });
  
    return formattedMessages;
  }