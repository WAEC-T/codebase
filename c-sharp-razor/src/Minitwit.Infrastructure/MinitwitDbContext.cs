using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;


namespace Minitwit.Infrastructure;

/// <summary>
///  EF Core will use the properties of the entities to create and control
///  the database, without having the application directly interact with the database.
/// </summary>

public sealed class MinitwitDbContext : IdentityDbContext<Author, IdentityRole<Guid>, Guid>
{
    public DbSet<Cheep> Cheeps { get; set; } = null!;

    public DbSet<Follow> Follows { get; set; } = null!;

    public DbSet<Reaction> Reactions { get; set; } = null!;

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
            string host = Environment.GetEnvironmentVariable("POSTGRES_HOST") ?? "192.168.178.41";
            string port = Environment.GetEnvironmentVariable("POSTGRES_PORT") ?? "5432";
            string database = Environment.GetEnvironmentVariable("POSTGRES_DB") ?? "waect";

            // Construct PostgreSQL connection string
            var connectionString = $"Host={host};Port={port};Database={database};Username={username};Password={password}";
            optionsBuilder.UseNpgsql(connectionString, 
                b => b.MigrationsAssembly("Minitwit.Infrastructure"));
        }
    }
    
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        // // Author entity
        // modelBuilder.Entity<Author>(entity =>
        // {
        //     modelBuilder
        //         .Entity<IdentityUserLogin<Guid>>()
        //         .HasKey(p => new { p.LoginProvider, p.ProviderKey });
        //     modelBuilder
        //         .Entity<IdentityUserLogin<Guid>>()
        //         .HasIndex(p => new { p.LoginProvider, p.ProviderKey })
        //         .IsUnique();
        //     modelBuilder.Entity<IdentityUserRole<Guid>>().HasKey(p => new { p.UserId, p.RoleId });
        //     modelBuilder
        //         .Entity<IdentityUserToken<Guid>>()
        //         .HasKey(p => new
        //         {
        //             p.UserId,
        //             p.LoginProvider,
        //             p.Name
        //         });
        //
        //     entity.Property(e => e.Id);
        // });
        //
        // modelBuilder.Entity<Follow>(entity =>
        // {
        //     entity.HasKey(f => new { f.FollowingAuthorId, f.FollowedAuthorId });
        // });
        //
        // // Cheep entity
        // modelBuilder.Entity<Cheep>(entity =>
        // {
        //     entity.HasKey(e => e.CheepId);
        //     entity.Property(e => e.AuthorId).IsRequired();
        //     entity.Property(e => e.Text).IsRequired();
        //     entity.Property(e => e.TimeStamp).IsRequired();
        // });
        //
        // modelBuilder.Entity<Reaction>().Property(m => m.ReactionType).HasConversion<string>();
        //
        // modelBuilder.Entity<Reaction>(entity =>
        // {
        //     entity.HasKey(r => new { r.CheepId, r.AuthorId });
        // });
        //
        // modelBuilder.Entity<IdentityUserLogin<Guid>>().HasKey(e => e.UserId);
        // modelBuilder.Entity<IdentityUserRole<Guid>>().HasKey(e => e.RoleId);
        // modelBuilder.Entity<IdentityUserToken<Guid>>().HasKey(e => e.UserId);

        base.OnModelCreating(modelBuilder); // Ensure the base configuration is applied

        // Configure Author entity
        modelBuilder.Entity<Author>(entity =>
        {
            entity.ToTable("users");
            entity.HasKey(a => a.Id);
            entity.Property(a => a.Id).HasColumnName("user_id");
            entity.Property(a => a.UserName).HasColumnName("user_name").HasMaxLength(50).IsRequired();;
            entity.Property(a => a.Email).HasColumnName("email").HasMaxLength(50).IsRequired();
        });

        // Configure Follow entity
        modelBuilder.Entity<Follow>(entity =>
        {
            entity.ToTable("followers");
            entity.HasKey(f => new { f.FollowingAuthorId, f.FollowedAuthorId });
            
            // Mapping columns
            entity.Property(f => f.FollowedAuthorId).HasColumnName("who_id");
            entity.Property(f => f.FollowedAuthorId).HasColumnName("whom_id");
            
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
                .HasColumnType("timestamp") // Correct type for PostgreSQL
                .HasDefaultValueSql("CURRENT_TIMESTAMP"); // Default value for timestamp
        });

        // Configure Reaction entity
        modelBuilder.Entity<Reaction>(entity =>
        {
            entity.ToTable("reactions"); // Update this if you have a different table name
            entity.HasKey(r => new { r.CheepId, r.AuthorId });
            
            entity.Property(r => r.CheepId).HasColumnName("cheep_id");
            entity.Property(r => r.AuthorId).HasColumnName("author_id");
            entity.Property(r => r.ReactionType).HasColumnName("reaction_type");
        });

        // Ensure indexes
        modelBuilder.Entity<Author>()
            .HasIndex(a => a.UserName)
            .IsUnique();

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
