using System.Configuration;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;


namespace Minitwit.Infrastructure;

/// <summary>
///  EF Core will use the properties of the entities to create and control
///  the database, without having the application directly interact with the database.
/// </summary>

public sealed class MinitwitDbContext : IdentityDbContext<Author, IdentityRole<int>, int>
{
    public DbSet<Message> Messages { get; set; } = null!;

    public DbSet<Follow> Follows { get; set; } = null!;

    public DbSet<Latest> LatestEntries { get; set; } = null!;
    
    public DbSet<Author> Authors { get; set; } = null!;
    
    public MinitwitDbContext(DbContextOptions<MinitwitDbContext> dbContextOptions)
        : base(dbContextOptions)
    {
        ChangeTracker.LazyLoadingEnabled = false;
    }

    protected override void OnModelCreating(ModelBuilder builder)
    {

        base.OnModelCreating(builder); // Ensure the base configuration is applied

        builder.Entity<Author>()
            .Ignore(u => u.AccessFailedCount)
            .Ignore(u => u.EmailConfirmed)
            .Ignore(u => u.LockoutEnabled)
            .Ignore(u => u.LockoutEnd)
            .Ignore(u => u.PhoneNumber)
            .Ignore(u => u.PhoneNumberConfirmed)
            .Ignore(u => u.SecurityStamp)
            .Ignore(u => u.TwoFactorEnabled);

        builder.Entity<Author>(entity =>
        {
            entity.ToTable("users");
            entity.HasKey(a => a.Id);
            entity.Property(a => a.Id).HasColumnName("user_id");
            entity.Property(a => a.UserName).HasColumnName("username").HasMaxLength(50).IsRequired();
            entity.Property(a => a.Email).HasColumnName("email").HasMaxLength(50).IsRequired();
            entity.Property(a => a.PasswordHash).HasColumnName("pw_hash").HasMaxLength(256).IsRequired();
            entity.Property(a => a.NormalizedUserName)
                .HasColumnName("normalized_username")
                .HasMaxLength(50).IsRequired();
            entity.Property(a => a.NormalizedEmail)
                .HasColumnName("normalized_email")
                .HasMaxLength(50).IsRequired();
            entity.Property(a => a.ConcurrencyStamp) 
                .HasColumnName("concurrency_stamp")
                .IsConcurrencyToken();
            entity.Property(a => a.SecurityStamp)
                .HasColumnName("security_stamp")
                .IsConcurrencyToken();
        });

        builder.Entity<Follow>(entity =>
        {
            entity.ToTable("followers");
            entity.HasKey(f => new { f.FollowingAuthorId, f.FollowedAuthorId });
            
            // Mapping columns
            entity.Property(f => f.FollowingAuthorId).HasColumnName("who_id");
            entity.Property(f => f.FollowedAuthorId).HasColumnName("whom_id");
            
            entity.HasOne<Author>().WithMany().HasForeignKey(f => f.FollowingAuthorId);
            entity.HasOne<Author>().WithMany().HasForeignKey(f => f.FollowedAuthorId);
        });

        builder.Entity<Message>(entity =>
        {
            entity.ToTable("messages");
            entity.HasKey(e => e.MessageId);
            
            entity.Property(e => e.MessageId).HasColumnName("message_id");
            entity.Property(e => e.AuthorId).HasColumnName("author_id");
            entity.Property(e => e.Text).HasColumnName("text").IsRequired();
            entity.Property(e => e.TimeStamp)
                .HasColumnName("pub_date")
                .HasColumnType("timestamp");
            entity.Property(e => e.Flagged).HasColumnName("flagged");
        });

        builder.Entity<Message>()
            .HasIndex(c => c.TimeStamp);

        builder.Entity<Latest>(entity =>
        {
            entity.ToTable("latest");

            entity.HasKey(e => e.Id);

            entity.Property(e => e.Id)
                .HasColumnName("id");

            entity.Property(e => e.Value)
                .HasColumnName("value");
        });
    }

    public async Task RemoveDuplicateUserLogins()
    {
        // Fetch all user logins
        var userLogins = await Set<IdentityUserLogin<Guid>>().ToListAsync();

        // Group by LoginProvider and ProviderKey
        var groupedUserLogins = userLogins.GroupBy(l => new { l.LoginProvider, l.ProviderKey });

        // For each group, keep only one record and mark the others for deletion
        foreach (var group in groupedUserLogins)
        {
            var userLoginsToDelete = group.Skip(1).ToList();
            Set<IdentityUserLogin<Guid>>().RemoveRange(userLoginsToDelete);
        }

        // Save changes to the database
        await SaveChangesAsync();
    }
}
