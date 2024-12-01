using FluentValidation;
using Minitwit.Core.Entities;

namespace Minitwit.Infrastructure.Repository;

public class MessageCreateValidator : AbstractValidator<CreateMessage>
{
    public MessageCreateValidator()
    {
        // @TODO Check that these values are correct:
        RuleFor(x => x.Text)
            .NotEmpty();
    }
}
