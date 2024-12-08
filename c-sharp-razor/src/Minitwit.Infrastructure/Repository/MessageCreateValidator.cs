using FluentValidation;
using Minitwit.Core.Entities;

namespace Minitwit.Infrastructure.Repository;

public class MessageCreateValidator : AbstractValidator<CreateMessage>
{
    public MessageCreateValidator()
    {
        RuleFor(x => x.Text)
            .NotEmpty();
    }
}
