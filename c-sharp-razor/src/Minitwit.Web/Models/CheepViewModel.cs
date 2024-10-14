using Minitwit.Core.Entities;

namespace Minitwit.Web.Models;

public record CheepViewModel(
    int CheepId,
    UserModel User,
    string Message,
    string Timestamp
);
