using Microsoft.AspNetCore.HttpLogging;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Infrastructure;
using Npgsql;
using DotNetEnv;

namespace Minitwit.Web;

static class ProgramOptions
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
        string? databaseUrl = Environment.GetEnvironmentVariable("DATABASE_URL");
        
        if (string.IsNullOrEmpty(databaseUrl))
        {
            throw new InvalidOperationException("DATABASE_URL environment variable is not set.");
        }
        
        try
        {
            // Parse DATABASE_URL into Npgsql-compatible format
            var connectionString = ConvertDatabaseUrlToConnectionString(databaseUrl);
        
            using var conn = new NpgsqlConnection(connectionString);

            builder.Services.AddDbContext<MinitwitDbContext>(options =>
            {
                options.UseNpgsql(connectionString);
            });
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Connection failed: {ex.Message}");
        }
    }
    
    private static string ConvertDatabaseUrlToConnectionString(string databaseUrl)
    {
        var uri = new Uri(databaseUrl);
        
        // Extract from URI
        var host = uri.Host;
        var port = uri.Port;
        var username = uri.UserInfo.Split(':')[0];
        var password = uri.UserInfo.Split(':')[1];
        var database = uri.AbsolutePath.Trim('/');
            
        // Build connection string
        return $"Host={host};Port={port};Username={username};Password={password};Database={database}";
    }
}
