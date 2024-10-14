using FluentValidation;
using Microsoft.EntityFrameworkCore;
using Minitwit.Core.Entities;
using Minitwit.Core.Repository;
using Minitwit.Infrastructure;
using Minitwit.Infrastructure.Repository;
using Minitwit.Web;

public class Program
{
    public static void Main(string[] args)
    {
        var builder = WebApplication.CreateBuilder(args);

        // Add services to the container.
        builder.Services.AddRazorPages();

        ProgramOptions.AddProgramOptions(builder);
        ProgramOptions.AddIdendity(builder);
        ProgramOptions.AddDatabase(builder);

        //API Controllers
        builder
            .Services.AddControllers()
            .AddJsonOptions(options =>
            {
                options.JsonSerializerOptions.PropertyNameCaseInsensitive = true;
                options.JsonSerializerOptions.PropertyNamingPolicy = null;
                options.JsonSerializerOptions.IgnoreNullValues = true;
            });

        // Dependency Injection
        builder.Services.AddScoped<IAuthorRepository, AuthorRepository>();
        builder.Services.AddScoped<IValidator<CreateCheep>, CheepCreateValidator>();
        builder.Services.AddScoped<ICheepRepository, CheepRepository>();
        builder.Services.AddScoped<ICheepService, MinitwitService>();
        builder.Services.AddScoped<IFollowRepository, FollowRepository>();
        builder.Services.AddScoped<ILatestRepository, LatestRepository>();

        builder.Services.AddSession(options =>
        {
            options.Cookie.Name = ".Minitwit.Web.Session";
            options.IdleTimeout = TimeSpan.FromMinutes(10);
            options.Cookie.HttpOnly = false;
            options.Cookie.IsEssential = true;
        });

        var app = builder.Build();

        // Apply migrations
        using (var scope = app.Services.CreateScope())
        {
            var services = scope.ServiceProvider;
            var context = services.GetRequiredService<MinitwitDbContext>();

            try
            {
                context.Database.Migrate();
                Console.WriteLine("Database migration applied successfully.");
            }
            catch (Exception e)
            {
                Console.WriteLine(e);
            }
        }

        if (!app.Environment.IsDevelopment())
        {
            app.UseExceptionHandler("/Error");
            app.UseHsts();
        }

        app.UseHttpsRedirection();
        app.UseStaticFiles();
        app.UseRouting();
        app.UseAuthentication();
        app.UseAuthorization();
        app.UseSession();
        app.MapControllers();
        app.MapRazorPages();

        app.Run();
    }
}
