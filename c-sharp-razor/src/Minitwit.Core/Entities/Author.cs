using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;

namespace Minitwit.Core.Entities;

/// <summary>
/// This class represents the user of the MiniTwit application.
/// All authenticated users are represented by an instance of this class.
/// Having the ability to create Messages, react to Messages,
/// and follow other users.
/// </summary>

[Index(nameof(Id), IsUnique = true)]
public class Author : IdentityUser<int>
{
    [StringLength(50)]
    [Required]
    public override required string UserName { get; set; }

    [StringLength(50)]
    [Required]
    public override required string Email { get; set; }
    
    [StringLength(50)]
    public override string? NormalizedEmail { get; set; } 
}
