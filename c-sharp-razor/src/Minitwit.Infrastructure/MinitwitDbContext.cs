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
    public DbSet<Cheep> Cheeps { get; set; } = null!;

    public DbSet<Follow> Follows { get; set; } = null!;
    
    public MinitwitDbContext(DbContextOptions<MinitwitDbContext> dbContextOptions)
        : base(dbContextOptions)
    {
        ChangeTracker.LazyLoadingEnabled = false;
    }
    
    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        if (!optionsBuilder.IsConfigured)
        {   
            // Get PostgreSQL environment variables
            string username = Environment.GetEnvironmentVariable("POSTGRES_USER") ?? "user";
            string password = Environment.GetEnvironmentVariable("POSTGRES_PASSWORD") ?? "pass";
            string host = Environment.GetEnvironmentVariable("POSTGRES_HOST") ?? "192.168.8.175"; // 
            string port = Environment.GetEnvironmentVariable("POSTGRES_PORT") ?? "5432";
            string database = Environment.GetEnvironmentVariable("POSTGRES_DB") ?? "waect";

            // Construct PostgreSQL connection string
            var connectionString = $"Host={host};Port={port};Database={database};Username={username};Password={password}";
            Console.WriteLine($"Connecting to PostgreSQL with connection string: {connectionString}");
            optionsBuilder.UseNpgsql(connectionString, 
                b => b.MigrationsAssembly("Minitwit.Infrastructure"));
        }
    }
    
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {   
        
        
        base.OnModelCreating(modelBuilder); // Ensure the base configuration is applied
        

        modelBuilder.Entity<Author>()
            .Ignore(u => u.AccessFailedCount)
            .Ignore(u => u.EmailConfirmed)
            .Ignore(u => u.LockoutEnabled)
            .Ignore(u => u.LockoutEnd)
            .Ignore(u => u.PhoneNumber)
            .Ignore(u => u.PhoneNumberConfirmed)
            .Ignore(u => u.SecurityStamp)
            .Ignore(u => u.TwoFactorEnabled);
        
        // Configure Author entity
        modelBuilder.Entity<Author>(entity =>
        {
            entity.ToTable("users");
            entity.HasKey(a => a.Id);
            entity.Property(a => a.Id).HasColumnName("user_id");
            entity.Property(a => a.UserName).HasColumnName("username").HasMaxLength(50).IsRequired();
            entity.Property(a => a.NormalizedUserName).HasColumnName("normalized_username").HasMaxLength(50).IsRequired();
            entity.Property(a => a.Email).HasColumnName("email").HasMaxLength(50).IsRequired();
            entity.Property(a => a.PasswordHash).HasColumnName("pw_hash").HasMaxLength(256).IsRequired();
            entity.Property(a => a.ConcurrencyStamp) //  changes to a user's profile or sensitive data (e.g., passwords) are safe 
                .HasColumnName("concurrency_stamp") 
                .IsConcurrencyToken();
            entity.Property(a => a.SecurityStamp) // cookies or refresh tokens
                .HasColumnName("security_stamp")
                .IsConcurrencyToken();
            
            entity.Ignore(a => a.NormalizedEmail);
        });

        // Configure Follow entity
        modelBuilder.Entity<Follow>(entity =>
        {
            entity.ToTable("followers");
            entity.HasKey(f => new { f.FollowingAuthorId, f.FollowedAuthorId });
            
            // Mapping columns
            entity.Property(f => f.FollowedAuthorId).HasColumnName("who_id");
            entity.Property(f => f.FollowingAuthorId).HasColumnName("whom_id");
            
            entity.HasOne<Author>().WithMany().HasForeignKey(f => f.FollowingAuthorId);
            entity.HasOne<Author>().WithMany().HasForeignKey(f => f.FollowedAuthorId);
        });

        // Configure Cheep entity
        modelBuilder.Entity<Cheep>(entity =>
        {
            entity.ToTable("messages");
            entity.HasKey(e => e.CheepId);
            
            entity.Property(e => e.CheepId).HasColumnName("message_id");
            entity.Property(e => e.AuthorId).HasColumnName("author_id");
            entity.Property(e => e.Text).HasColumnName("text").IsRequired();
            entity.Property(e => e.TimeStamp)
                .HasColumnName("pub_date")
                .HasColumnType("text")
                .HasConversion(
                    v => v.ToString("yyyy-MM-dd HH:mm:ss"),  // Convert DateTime to string when saving
                    v => DateTime.Parse(v)                   // Convert string to DateTime when reading
                );
        });

        modelBuilder.Entity<Cheep>()
            .HasIndex(c => c.TimeStamp);
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
