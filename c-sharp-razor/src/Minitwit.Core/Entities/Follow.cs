using System.ComponentModel.DataAnnotations;

namespace Minitwit.Core.Entities;

/// <summary>
/// This class represents a follow between two authors.
/// A follow is a one way relationship between two authors.
/// The FollowingAuthor is the author who is following the FollowedAuthor.
/// The FollowedAuthor is the author who is being followed by the FollowingAuthor.
/// The FollowingAuthorId is the id of the FollowingAuthor.
/// The FollowedAuthorId is the id of the FollowedAuthor.
/// </summary>

public class Follow
{
    [Required]
    public int FollowingAuthorId { get; set; }

    [Required]
    public int FollowedAuthorId { get; set; }
}
