using Microsoft.AspNetCore.HttpLogging;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Infrastructure;
using Npgsql;
using DotNetEnv;

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
        string dbUser = Environment.GetEnvironmentVariable("POSTGRES_USER") ?? "user";
        string dbPassword = Environment.GetEnvironmentVariable("POSTGRES_PASSWORD") ?? "pass";
        string dbHost = Environment.GetEnvironmentVariable("POSTGRES_HOST") ?? "local_postgres";
        string dbPort = Environment.GetEnvironmentVariable("POSTGRES_PORT") ?? "5432";
        string dbName = Environment.GetEnvironmentVariable("POSTGRES_DB") ?? "waect";
        
        var connectionString = $"Host={dbHost};Port={dbPort};Database={dbName};Username={dbUser};Password={dbPassword}";
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

        Console.WriteLine(
            $"Connection string: Host={dbHost}; Port={dbPort}; Database={dbName}; Username={dbUser}; Password={dbPassword}");
    }
}
