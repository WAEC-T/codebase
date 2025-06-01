using System.ComponentModel.DataAnnotations;
using Microsoft.EntityFrameworkCore;

namespace Minitwit.Core.Entities;

/// <summary>
/// This class represents Messages, created by the users of the MiniTwit application.
/// Messages are small messages,5 to 160 characters in length.
/// It's the only way for users to communicate with each other.
/// Messages hold the reactions given to them.
/// </summary>

[Index(nameof(MessageId), IsUnique = true)]
public class Message
{
    [Required]
    public int MessageId { get; set; }

    [Required]
    public int AuthorId { get; set; }

    [Required]
    public required string Text { get; set; }

    [Required]
    public DateTime TimeStamp { get; set; }
    
    [Required]
    public int Flagged { get; set; } 
}
