using Minitwit.Core.Entities;

namespace Minitwit.Web.Models;

public record MessageViewModel(
    int MessageId,
    UserModel User,
    string Message,
    DateTime Timestamp,
    int Flagged
);
