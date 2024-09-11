using Microsoft.AspNetCore.HttpLogging;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Infrastructure;
using Npgsql;

namespace Minitwit.Web;

public class ProgramOptions
{
    public static void AddProgramOptions(WebApplicationBuilder builder)
    {
        builder.Services.AddHttpLogging(logging =>
        {
            logging.LoggingFields = HttpLoggingFields.All;
            logging.MediaTypeOptions.AddText("application/javascript");
            logging.RequestBodyLogLimit = 0;
            logging.ResponseBodyLogLimit = 0;
            logging.CombineLogs = true;
        });
    }

    public static void AddIdendity(WebApplicationBuilder builder)
    {
        builder.Services.Configure<IdentityOptions>(options =>
        {
            options.Password.RequireDigit = false;
            options.Password.RequiredLength = 1;
            options.Password.RequireLowercase = false;
            options.Password.RequireUppercase = false;
            options.Password.RequireNonAlphanumeric = false;
            options.User.AllowedUserNameCharacters =
                "zxcvbnmasdfghjklqwertyuiopZXCVBNMASDFGHJKLQWERTYUIOP1234567890 @";
        });

        builder
            .Services.AddDefaultIdentity<Author>()
            .AddRoles<IdentityRole<int>>()
            .AddEntityFrameworkStores<MinitwitDbContext>();

        builder.Services.AddAuthentication().AddCookie();
    }

    public static void AddDatabase(WebApplicationBuilder builder)
    {
        string environmentName = Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT");

        // Set up the database path
        if (environmentName != null && environmentName.Equals("Development"))
        {
            string currentDirectory = Directory.GetCurrentDirectory();
            string dbPath;

            if (
                Directory.Exists(
                    Path.Combine(currentDirectory, "..", "Minitwit.Infrastructure", "data")
                )
            )
            {
                dbPath = Path.Combine(
                    currentDirectory,
                    "..",
                    "Minitwit.Infrastructure",
                    "data",
                    "MinitwitDBContext.db"
                ); //Build directory
            }
            else
            {
                dbPath = Path.Combine(currentDirectory, "data", "MinitwitDBContext.db"); //Publish directory
            }
            builder.Services.AddDbContext<MinitwitDbContext>(options =>
            {
                options.UseSqlite($"Data Source={dbPath}");
            });
        }
        else
        {
            // Get PostgreSQL environment variables
            string username = Environment.GetEnvironmentVariable("POSTGRES_USER") ?? "user";
            string password = Environment.GetEnvironmentVariable("POSTGRES_PASSWORD") ?? "pass";
            string host = Environment.GetEnvironmentVariable("POSTGRES_HOST") ?? "192.168.8.175";
            string port = Environment.GetEnvironmentVariable("POSTGRES_PORT") ?? "5432";
            string database = Environment.GetEnvironmentVariable("POSTGRES_DB") ?? "waect";

            // Construct PostgreSQL connection string
            var connectionString = $"Host={host};Port={port};Database={database};Username={username};Password={password}";
            Console.WriteLine($"Testing connection with: {connectionString}");

            try
            {
                using var conn = new NpgsqlConnection(connectionString);
                conn.Open();
                Console.WriteLine("Connection successful!");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Connection failed: {ex.Message}");
            }
            
            builder.Services.AddDbContext<MinitwitDbContext>(options =>
            {
                options.UseNpgsql(connectionString);
            });

            Console.WriteLine($"Connection string: Host={host}; Port={port}; Database={database}; Username={username}; Password={password}");
        }
    }
}
