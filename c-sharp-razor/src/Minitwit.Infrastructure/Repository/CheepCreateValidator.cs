using FluentValidation;
using Minitwit.Core.Entities;

namespace Minitwit.Infrastructure.Repository;

public class MessageCreateValidator : AbstractValidator<CreateMessage>
{
    public MessageCreateValidator()
    {
        // @TODO Check that these values are correct:
        RuleFor(x => x.Text)
            .NotEmpty()
            .MaximumLength(160)
            .MinimumLength(5)
            .WithMessage("The Message must be between 5 and 160 characters.(MessageCreateValidator)");
    }
}
