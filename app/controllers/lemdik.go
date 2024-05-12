package controllers

import "github.com/gofiber/fiber/v2"


func RegisterLemdik(c *fiber.Ctx)error{

	//Pake Role Super admin/ admin pusat 





	return c.JSON(fiber.Map{
		"Pesan":"Sukses Membuat Lemdik",

	})
}

func LoginLemdik(c *fiber.Ctx)error{

	//Pake Role Super admin/ admin pusat 





	return c.JSON(fiber.Map{
		"t":"",

	})
}

func GetLemdik(c *fiber.Ctx)error{

	//id:= c.Params("id")


	return c.JSON(fiber.Map{
		"Pesan":"",
		"data":"",
	})
}

func UpdateLemdik(c *fiber.Ctx)error{




	return c.JSON(fiber.Map{
		"Pesan":"",

	})
}

func DeleteLemdik(c *fiber.Ctx)error{


	return nil
}




